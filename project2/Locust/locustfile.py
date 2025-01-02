import time, json, random
from locust import HttpUser, task

class PostCoursesTest(HttpUser):

  @task
  def post_grades(self):
    courses = []
    with open("courses.json", "r") as file:
      courses = json.load(file)
    
    for _ in range(len(courses)):
      response = self.client.post("http://localhost:8000/course", json=courses.pop(), name="/coursePost")
      print(response.text)
      print(response.status_code)
      time.sleep(0.3)

  def on_start(self):
    print("task started")

  def on_stop(self):
    print("task end")