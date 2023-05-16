allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/database.Tiltfile', 'redis', 'cockroach', 'migrations')
load('./deploy/tilt/app.Tiltfile', 'apps')

hostname = 'ponglehub.co.uk'

k8s_yaml(namespace_yaml('book-planner'))

operator('db')
cockroach('cockroach', '256Mi')

migrations('src/cmd/apis', 'cockroach')
apps('src/cmd/apis', hostname)
apps('src/cmd/pages', hostname)