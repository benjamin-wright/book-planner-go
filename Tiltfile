allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/database.Tiltfile', 'redis', 'cockroach')
load('./deploy/tilt/app.Tiltfile', 'secure_api', 'insecure_api', 'internal_api')
load('./deploy/tilt/app.Tiltfile', 'secure_page', 'insecure_page')

basepath = 'http://localhost'

k8s_yaml(namespace_yaml('book-planner'))

operator('db')
redis('auth-redis', '128Mi')

insecure_api(basepath, 'apis', 'auth', ['env.LOGIN_URL=%s/login' % basepath])
insecure_page(basepath, 'pages', 'login')

secure_page(basepath, 'pages', 'home')