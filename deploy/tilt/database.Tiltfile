def redis(name, storage):
    db(name, 'redis', storage)

def cockroach(name, storage):
    db(name, 'cockroach', storage)

def db(name, db_type, storage):
    k8s_yaml(helm(
        'deploy/helm/database',
        name='db-%s-%s' % (db_type, name),
        namespace='book-planner',
        set=[
            'name=%s' % name,
            'type=%s' % db_type,
            'storage=%s' % storage,
        ],
    ))

    k8s_resource(
        new_name=name,
        objects=['%s:%sdb' % (name, db_type)],
        extra_pod_selectors=[{'app': name}],
    )