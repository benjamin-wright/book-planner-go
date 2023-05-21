load('ext://helm_remote', 'helm_remote')

def nats(name, labels=[]):
    helm_remote(
        chart = 'nats',
        repo_name = 'nats',
        repo_url = 'https://nats-io.github.io/k8s/helm/charts/',
        release_name = name,
        namespace = 'book-planner',
        set=[
            'podDisruptionBudget.enabled=false',
        ]
    )

    k8s_resource(
        '%s-nats' % name,
        labels=['infra']+labels
    )

    k8s_resource(
        '%s-nats-box' % name,
        labels=['infra']+labels
    )