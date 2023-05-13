def app(path, name, base_url, labels=[]):
    custom_build(
        name,
        'just build %s $EXPECTED_REF' % path,
        [ path ],
        ignore = [
            'dist/*',
            '**/*_test.go'
        ]
    )

    # local_resource(
    #     '%s-test' % fullname,
    #     'just test src/%s/%s' % (path, name),
    #     auto_init = False,
    #     trigger_mode = TRIGGER_MODE_MANUAL
    # )

    # local_resource(
    #     '%s-int-test' % fullname,
    #     'just int-test src/%s/%s' % (path, name),
    #     auto_init = False,
    #     trigger_mode = TRIGGER_MODE_MANUAL
    # )

    k8s_yaml(helm(
        'deploy/helm/app',
        name=name,
        namespace='book-planner',
        values=['%s/values.yaml' % path],
        set=[
            "name=%s" % name,
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