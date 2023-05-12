allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./deploy/tilt/operator.Tiltfile', 'operator')
load('./deploy/tilt/database.Tiltfile', 'redis', 'cockroach')
load('./deploy/tilt/app.Tiltfile', 'app')

base_url = 'http://localhost'

k8s_yaml(namespace_yaml('book-planner'))

operator('db')

redis('redis', '128Mi')

app('src/cmd/apis/auth', 'apis-auth', base_url)
app('src/cmd/pages/login', 'pages-login', base_url)
app('src/cmd/pages/register', 'pages-register', base_url)

app('src/cmd/pages/home', 'pages-home', base_url)