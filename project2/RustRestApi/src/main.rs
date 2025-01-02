use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder, Result};
use serde::Deserialize;
use std::env;

#[derive(Deserialize)]
struct Course {
    curso: String,
    facultad: String,
    carrera: String,
    region: String
}

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello sopes!")
}

#[post("/course")]
async fn course(course: web::Json<Course>) -> Result<String> {
    Ok(format!("actix, course received: {} {} {} {}", course.curso, course.facultad, course.carrera, course.region))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let host = env::var("RUST_SERVER_HOST").unwrap();
    let port = env::var("RUST_SERVER_PORT").unwrap();

    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(course)    
    })
    .bind((host, port.parse().unwrap()))?
    .run()
    .await
}