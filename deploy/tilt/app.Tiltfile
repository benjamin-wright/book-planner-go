def app(path, name, secure):
    custom_build(
        '%s-%s' % (path, name),
        'just build src/%s/%s $EXPECTED_REF' % (path, name),
        [
            'src/%s/%s' % (path, name)
        ],
        ignore = [
            'dist/*',
            '**/*_test.go'
        ]
    )

    local_resource(
        '%s-%s-test' % (path, name),
        'just test src/%s/%s' % (path, name),
        deps = [name],
        auto_init = False,
        trigger_mode = TRIGGER_MODE_MANUAL
    )

    local_resource(
        '%s-%s-int-test' % (path, name),
        'just int-test src/%s/%s' % (path, name),
        deps = [name],
        auto_init = False,
        trigger_mode = TRIGGER_MODE_MANUAL
    )

    k8s_yaml(helm(
        'deploy/helm/app',
        name='%s-%s' % (path, name),
        namespace='book-planner',
        set=[
            'name=%s-%s' % (path, name),
            'image=%s-%s' % (path, name),
            'path=/%s' % (name),
            'secure=%s' % secure
        ]
    ))

    k8s_resource(
        '%s-%s' % (path, name),
        auto_init = True,
        trigger_mode = TRIGGER_MODE_MANUAL
    )