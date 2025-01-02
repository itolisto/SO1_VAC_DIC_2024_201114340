import random, io, json

class GradesJsonGenerator:

  def _getCarnet(self):
    years = random.choice([2023, 2024])
    number = random.randint(00000, 99999)
    return f"{years}{number}"

  def _getName(self):
    name = random.choice(["Pablo", "Fernando", "Mario", "Maria", "Karla", "Yahaira"])
    lastName = random.choice(["Torres", "Puac", "Boch", "Gonzalez", "Lopez"])
    return f"{name} {lastName}"

  def _getCourse(self):
    return random.choice(["SO1", "BD1", "LFP", "SA", "AYD1"])

  def _getGrade(self):
    return random.choice([10, 20, 30, 40, 50, 60, 70, 80, 90, 100])

  def _getSemester(self):
    return random.choice(["1S", "2S"])

  def _getYear(self):
    return random.randint(2009, 2024)

  def __init__(self, fileName, entries):
    self.grades = []

    for x in range(entries - 1):
      gradeEntry = { 'carnet': self._getCarnet(), 'nombre': self._getName(), 'curso': self._getCourse(), 'nota': self._getGrade(), 'semestre': self._getSemester(), 'año': self._getYear() }
      self.grades.append(gradeEntry)
    
    with open(f"./{fileName}.json", 'w') as file:
      json.dump(self.grades, file, indent = 2, ensure_ascii=False)

generator = GradesJsonGenerator("grades", 100)        