def operator(name):
    dirs = listdir('src/cmd/operators/%s/crds' % name)

    custom_build(
        name,
        'just build src/cmd/operators/%s $EXPECTED_REF' % name,
        [
            'src/cmd/operators/%s' % name
        ],
        ignore = [
            'dist/*',
            '**/*_test.go'
        ]
    )

    # local_resource(
    #     '%s-test' % name,
    #     'just test src/cmd/operators/%s' % name,
    #     deps = [name],
    #     auto_init = False,
    #     trigger_mode = TRIGGER_MODE_MANUAL
    # )

    # local_resource(
    #     '%s-int-test' % name,
    #     'just int-test src/cmd/operators/%s' % name,
    #     deps = [name],
    #     auto_init = False,
    #     trigger_mode = TRIGGER_MODE_MANUAL
    # )

    k8s_yaml(dirs)

    k8s_yaml(helm(
        'deploy/helm/operator',
        name=name,
        namespace='book-planner',
        values=[
            'src/cmd/operators/%s/values.yaml' % name
        ],
        set=[
            'image=%s' % name
        ]
    ))

    k8s_resource(
        name,
        auto_init = True,
        trigger_mode = TRIGGER_MODE_MANUAL
    )