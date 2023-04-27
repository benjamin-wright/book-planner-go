allow_k8s_contexts(['book-planner'])
load('ext://namespace', 'namespace_yaml')
load('./src/operators/Tiltfile', 'operator')

k8s_yaml(namespace_yaml('book-planner'))

operator('db')