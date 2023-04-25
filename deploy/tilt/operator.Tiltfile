def operator(name, path):
    dirs = listdir(path+'/crds')

    k8s_yaml(dirs)

    k8s_yaml(helm(
        'deploy/operator',
        name=name,
        namespace='book-planner',
        values=[
            path+'/values.yaml'
        ],
        set=[
            'image=localhost:5000/%s' % name
        ]
    ))