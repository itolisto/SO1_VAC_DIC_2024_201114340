import express from 'express'
import mysqlx from '@mysql/xdevapi'

const app = express()

app.use(express.json())

const dbname = process.env.DB_NAME

var client = mysqlx.getClient(
  {
    host            : process.env.DB_HOST,
    user            : process.env.DB_USER,
    password        : process.env.DB_PASSWORD,
    port            : process.env.DB_PORT       // x protocol port, using 3306
  }, 
  { 
    pooling: { 
      enabled: true,
      maxIdleTime: 30000,
      maxSize: 5, 
      queueTimeout: 10000 
    } 
  }
)

app.listen(8080, () => { console.log(`app listening on port 8080, host:${process.env.DB_HOST}`) });

app.get('/', (req, res) => {
    res.send('Bishitus bebe')
})

app.post('/ram', (req, res, next) => {
  let ipAddress = req.header('x-forwarded-for')

  if (ipAddress == undefined) {
    ipAddress = req.socket.remoteAddress
  }

  if (ipAddress.startsWith("::ffff:")) ipAddress = ipAddress.substring(7)

  const body = req.body

  console.log(`inserting ram ${body}`)

  client
    .getSession()
    .then( session => {
      session.sql(`USE ${dbname}`).execute();

      session.sql('INSERT IGNORE INTO vm (ip) VALUES (?)').bind(ipAddress).execute()

      console.log(body.total_ram + " " + body.free_ram + " " + body.used_ram + " " + body.percentage_used + " " + ipAddress)
     
      session
        .sql('INSERT INTO ram (total_ram, free_ram, used_ram, percentage_used, ip) VALUES (?, ?, ?, ?, ?)')
        .bind(body.total_ram, body.free_ram, body.used_ram, body.percentage_used, ipAddress)
        .execute()
        .catch(function (err) {
          try {
            session.close()
            res.sendStatus(400)
            next(err)  // expressjs error handling
            console.log("ram not inserted")
            return session.close()  
          } catch (error) {
            console.log("ram catch catch" + error.message)
          }
          
        }).then( _ => {
          try {
            console.log("ram inserted")
            res.send('ram info inserted')  
            return session.close() 
          } catch (error) {
            console.log("ram inserted catch" + error.message)
          }
        })

    })
    .catch(function (err) {
      next(err)
      console.log('data base error: ' + err.message);
    })
})

app.post('/cpu', (req, res, next) => {
  let ipAddress = req.socket.remoteAddress

  if (ipAddress == undefined) {
    ipAddress = req.socket.remoteAddress
  }

  if (ipAddress.startsWith("::ffff:")) ipAddress = ipAddress.substring(7)

  const body = req.body

  console.log(`inserting cpu`)

    let free = 100 - parseFloat(body.percentage_used)
    if(free == undefined) free = 100.0 - parseFloat(body.percentage_used)
    

    client
    .getSession()
    .then( session => {
      session.sql(`USE ${dbname}`).execute();

      session.sql('INSERT IGNORE INTO vm (ip) VALUES (?)').bind(ipAddress).execute()

      session
        .sql('INSERT IGNORE INTO cpu (percentage_used, free, ip) VALUES (?, ?, ?)')
        .bind(body.percentage_used, free, ipAddress)
        .execute()
        .catch(function (err) {
          next(err)  // expressjs error handling
          console.log("cpu not inserted "+ err.message)
          res.sendStatus(400)
        })
        .then( result => {
          console.log('cpu row inserted')

          try {
            res.send('cpu inserted')  
          } catch (error) {
            
          }
        }) 

      console.log(`inserting processes`)     

      body.tasks?.forEach(task => {
        console.log(task.pid + " " + task.name + " " + task.state + " " + task.user + " " + task.ram + " " + task.father + " " + ipAddress)
        session
          .sql('INSERT INTO process (pid, name, state, puser, ram, father, ip) VALUES (?, ?, ?, ?, ?, ?, ?)')
          .bind(task.pid, task.name, task.state, task.user, task.ram, task.father, ipAddress)
          .execute()
          .catch(function (err) {
            next(err)  // expressjs error handling
            try {
              res.sendStatus(400)
              console.log("processes rows not inserted" + err.message)
            } catch (error) {
              
            }
          })
          .then( result => {
            console.log('process row inserted')

            try {
              res.send('process inserted')  
            } catch (error) {
              
            }
            
            session.close()
          }) 
      })

      if (body.tasks == undefined) {
        session.close()
      }
    })
  })