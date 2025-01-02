import random, io, json

class GradesJsonGenerator:

  def _getCurso(self):
    return random.choice(["SO1", "ANP", "IO2", "SA", "OLC2"])

  def _getFacultad(self):
    return random.choice(["Ingenieria", "Medicina", "Arquitectura", "Humanidades"])

  def _getCarrera(self):
    return random.choice(["Sistemas", "Industrial", "Diseño Grafico", "Arquitectura", "Pedadogia", "Arte"])

  def _getRegion(self):
    return random.choice(["METROPOLITANA", "NORTE", "NORORIENTAL", "SURORIENTAL", "CENTRAL", "SUROCCIDENTAL", "NOROCCIDENTAL", "PETEN"])

  def __init__(self, fileName, entries):
    self.grades = []

    for x in range(entries - 1):
      gradeEntry = { 'curso': self._getCurso(), 'facultad': self._getFacultad(), 'carrera': self._getCarrera(), 'region': self._getRegion() }
      self.grades.append(gradeEntry)
    
    with open(f"./{fileName}.json", 'w') as file:
      json.dump(self.grades, file, indent = 2, ensure_ascii=False)

generator = GradesJsonGenerator("courses", 100)        