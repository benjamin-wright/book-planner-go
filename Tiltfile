allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/database.Tiltfile', 'redis', 'cockroach', 'migrations')
load('./deploy/tilt/nats.Tiltfile', 'nats')
load('./deploy/tilt/app.Tiltfile', 'apps')
load('./deploy/tilt/events.Tiltfile', 'events')

hostname = 'ponglehub.co.uk'

k8s_yaml(namespace_yaml('book-planner'))

operator('db')

nats('events')
cockroach('cockroach', '256Mi')

migrations('src/cmd/apis', 'cockroach')
apps('src/cmd/apis', hostname)
events('src/cmd/events', 'nats://events-nats.book-planner.svc.cluster.local:4222')
apps('src/cmd/pages', hostname)