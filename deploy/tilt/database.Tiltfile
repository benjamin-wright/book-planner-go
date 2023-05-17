load('ext://helm_resource', 'helm_resource')

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
        labels=['infra'],
    )

def load_migrations(path):
    file = migration_name(path)
    [ database, index ] = file.split('-')
    migration = str(local(
        'cat %s' % path,
        echo_off=True,
        quiet=True,
    )).rstrip('\n').split('\n')

    return [
        '--set=migrations.%s.database=%s' % (file, database),
        '--set=migrations.%s.index=%s' % (file, index),
        '--set-file=migrations.%s.migration=%s' % (file, path),
    ]

def migration_name(path):
    return path.split('/')[-1].replace('.sql', '')

def migrations(path, db):
    dirs = str(local(
        'find %s -name migrations' % path,
        echo_off=True,
        quiet=True,
    )).rstrip('\n').split('\n')

    for d in dirs:
        files = str(local(
            'find %s -name *.sql' % d,
            echo_off=True,
            quiet=True,
        )).rstrip('\n').split('\n')

        flags = [ flag for file in files for flag in load_migrations(file) ]
        parent_dir = d.split('/')[-2]

        helm_resource(
            'mig-%s' % parent_dir,
            'deploy/helm/migrations',
            namespace='book-planner',
            flags=[
                '--set=deployment=%s' % db,
            ] + flags,
        )