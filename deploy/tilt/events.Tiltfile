def has_element(items, item):
    for i in items:
        if i == item:
            return True

    return False

def events(path, natsUrl):
    files = str(local(
        'find %s -name values.yaml' % path,
        echo_off=True,
        quiet=True,
    )).rstrip('\n').split('\n')

    names = list()

    for file in files:
        parts = file.split('/')
        values = read_yaml(file)

        name = values['name']
        labels = values.get('labels', [])
        new_path = '/'.join(parts[0:-1])

        if has_element(names, name):
            fail('duplicated name: %s' % name)
        else:
            names.append(values["name"])

        custom_build(
            name,
            'just build %s $EXPECTED_REF' % new_path,
            [ new_path ],
            ignore = [
                'dist/*',
                '**/*_test.go'
            ]
        )

        k8s_yaml(helm(
            'deploy/helm/events',
            name=name,
            namespace='book-planner',
            values=[file],
            set=[
                "image=%s" % name,
                "natsUrl=%s" % natsUrl,
            ]
        ))

        k8s_resource(
            name,
            auto_init = True,
            trigger_mode = TRIGGER_MODE_MANUAL,
            labels=labels,
        )
