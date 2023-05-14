def has_element(items, item):
    for i in items:
        if i == item:
            return True

    return False

def apps(path, base_url):
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

        # local_resource(
        #     '%s-test' % fullname,
        #     'just test src/%s/%s' % (new_path, name),
        #     auto_init = False,
        #     trigger_mode = TRIGGER_MODE_MANUAL
        # )

        # local_resource(
        #     '%s-int-test' % fullname,
        #     'just int-test src/%s/%s' % (new_path, name),
        #     auto_init = False,
        #     trigger_mode = TRIGGER_MODE_MANUAL
        # )

        k8s_yaml(helm(
            'deploy/helm/app',
            name=name,
            namespace='book-planner',
            values=[file],
            set=[
                "image=%s" % name,
                "baseUrl=%s" % base_url
            ]
        ))

        k8s_resource(
            name,
            auto_init = True,
            trigger_mode = TRIGGER_MODE_MANUAL,
            labels=labels,
        )
