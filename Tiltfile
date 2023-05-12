allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/app.Tiltfile', 'secure', 'insecure')

k8s_yaml(namespace_yaml('book-planner'))

operator('db')

insecure('apis', 'auth', ['env.LOGIN_URL=http://localhost/login'])

secure('pages', 'home', [])