allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/database.Tiltfile', 'redis', 'cockroach')
load('./deploy/tilt/app.Tiltfile', 'apps')

hostname = 'ponglehub.co.uk'

k8s_yaml(namespace_yaml('book-planner'))

operator('db')
redis('redis', '128Mi')

apps('src/cmd/apis', hostname)
apps('src/cmd/pages', hostname)