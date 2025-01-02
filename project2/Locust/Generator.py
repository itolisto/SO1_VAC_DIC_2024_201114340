import random

def getCarnet():
  years = random.choice([2023, 2024])
  number = random.randint(00000, 99999)
  return f"{years}{number}"

def getName():
  name = random.choice(["Pablo", "Fernando", "Mario", "Maria", "Karla", "Yahaira"])
  lastName = random.choice(["Torres", "Puac", "Boch", "Gonzalez", "Lopez"])
  return f"{name} {lastName}"

def getCourse():
  return random.choice(["SO1", "BD1", "LFP", "SA", "AYD1"])

def getGrade():
  return random.choice([10, 20, 30, 40, 50, 60, 70, 80, 90, 100])

def getSemester():
  return random.choice(["1S", "2S"])

def getYear():
  return random.randint(2009, 2024)