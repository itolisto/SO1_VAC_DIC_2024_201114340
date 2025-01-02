use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use redis::{AsyncCommands, Client, JsonAsyncCommands};
use serde::{Deserialize, Serialize};
// use serde_json;
use std::env;

#[derive(Deserialize, Serialize)]
struct Course {
    curso: String,
    facultad: String,
    carrera: String,
    region: String,
}

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello sopes!")
}

#[post("/course")]
async fn course(course: web::Json<Course>) -> impl Responder {
    let redis_url = format!(
        "redis://{}:{}/",
        env::var("RUST_REDIS_HOST").unwrap(),
        env::var("RUST_REDIS_PORT").unwrap()
    );

    let client = match Client::open(redis_url) {
        Ok(rclient) => rclient,
        Err(e) => {
            return HttpResponse::InternalServerError()
                .body(format!("Error connecting to Redis: {}", e))
        }
    };

    let mut con = match client.get_multiplexed_async_connection().await {
        Ok(connection) => connection,
        Err(e) => {
            return HttpResponse::InternalServerError()
                .body(format!("Error connecting to Redis client: {}", e))
        }
    };

    // let course_json = match serde_json::to_string(&course) {
    //     Ok(j) => j,
    //     Err(e) => return HttpResponse::InternalServerError().body(format!("Error parssing back to json: {}", e))
    // };

    if let Err(e) = con
        .json_set::<&str, &str, web::Json<Course>, ()>("assignacion:2", "$", &course)
        .await
    {
        return HttpResponse::InternalServerError()
            .body(format!("Error setting json to redis: {}", e));
    };

    // let result = match con.get::<&str, isize>("my_key").await {
    //     Ok(connection) => connection,
    //     Err(e) => return HttpResponse::InternalServerError().body(format!("Error getting key: {}", e))
    // };

    // println!("json in Rust server is: {}", result);

    HttpResponse::Ok().body("actix, course received")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let host = env::var("RUST_SERVER_HOST").unwrap();
    let port = env::var("RUST_SERVER_PORT").unwrap();

    HttpServer::new(|| App::new().service(hello).service(course))
        .bind((host, port.parse().unwrap()))?
        .run()
        .await
}
