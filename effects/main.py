import sys
import matplotlib.pyplot as plt
from scipy.spatial import Delaunay
from point_generators import generate_max_entropy_points, edge_points, generate_uniform_random_points
import numpy as np
import json

image = plt.imread(sys.argv[1])[:,:,:3]
n_points = int(sys.argv[2])

points = generate_uniform_random_points(image, n_points)


tri = Delaunay(points)

triangles = {
    'points': tri.points.astype(int).tolist(),
    'simplices': tri.simplices.tolist()
}

with open('triangles.json', 'w') as json_file:
  json.dump(triangles, json_file)