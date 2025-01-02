import time, json, random
from locust import HttpUser, task

class PostGradesTest(HttpUser):

  @task
  def post_grades(self):
    grades = []
    with open("grades.json", "r") as file:
      grades = json.load(file)
    
    for _ in range(len(grades)):
      response = self.client.post("http://localhost:8000/grade", json=grades.pop(), name="/gradePost")
      print(response.text)
      print(response.status_code)
      time.sleep(random.choice([1,2]))

  def on_start(self):
    print("task started")

  def on_stop(self):
    print("task end")