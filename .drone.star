# drone-npm build

def main(ctx):
  before = testing()

  stages = [
    linux('amd64'),
    linux('arm64'),
    linux('arm'),

  ]

  after = manifest() + gitter()

  for b in before:
    for s in stages:
      s['depends_on'].append(b['name'])

  for s in stages:
    for a in after:
      a['depends_on'].append(s['name'])

  return before + stages + after

def testing():
  return [{
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'testing',
    'platform': {
      'os': 'linux',
      'arch': 'amd64',
    },
    'steps': [
      {
        'name': 'vet',
        'image': 'golang:1.13',
        'pull': 'always',
        'commands': [
          'go version',
          'go vet ./...'
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go'
          }
        ]
      },
      {
        'name': 'test',
        'image': 'golang:1.13',
        'pull': 'always',
        'commands': [
          'go version',
          'go test -cover ./...'
        ],
        'volumes': [
          {
            'name': 'gopath',
            'path': '/go'
          }
        ]
      }
    ],
    'volumes': [
      {
        'name': 'gopath',
        'temp': {}
      }
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**'
      ]
    }
  }]

def linux(arch):
  return {
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'linux-%s' % arch,
    'platform': {
      'os': 'linux',
      'arch': arch,
    },
    'steps': [
      {
        'name': 'build-push',
        'image': 'golang:1.13',
        'pull': 'always',
        'environment': {
          'CGO_ENABLED': '0'
        },
        'commands': [
          'go version',
          'go build -v -ldflags "-X main.version=${DRONE_COMMIT_SHA:0:8}" -a -tags netgo -o release/linux/%s/drone-npm ./cmd/drone-npm' % arch,
        ],
        'when': {
          'event': {
            'exclude': [
              'tag'
            ]
          }
        }
      },
      {
        'name': 'build-tag',
        'image': 'golang:1.13',
        'pull': 'always',
        'environment': {
          'CGO_ENABLED': '0'
        },
        'commands': [
          'go version',
          'go build -v -ldflags "-X main.version=${DRONE_TAG##v}" -a -tags netgo -o release/linux/%s/drone-npm ./cmd/drone-npm' % arch,
        ],
        'when': {
          'event': [
            'tag'
          ]
        }
      },
      {
        'name': 'executable',
        'image': 'golang:1.13',
        'pull': 'always',
        'commands': [
          './release/linux/%s/drone-npm --help' % arch
        ]
      },
      {
        'name': 'dryrun',
        'image': 'plugins/docker',
        'pull': 'always',
        'settings': {
          'dry_run': True,
          'tags': 'linux-%s' % arch,
          'dockerfile': 'docker/Dockerfile.linux.%s' % arch,
          'repo': 'plugins/npm',
          'username': {
            'from_secret': 'docker_username'
          },
          'password': {
            'from_secret': 'docker_password'
          }
        },
        'when': {
          'event': [
            'pull_request'
          ]
        }
      },
      {
        'name': 'publish',
        'image': 'plugins/docker',
        'pull': 'always',
        'settings': {
          'auto_tag': True,
          'auto_tag_suffix': 'linux-%s' % arch,
          'dockerfile': 'docker/Dockerfile.linux.%s' % arch,
          'repo': 'plugins/npm',
          'username': {
            'from_secret': 'docker_username'
          },
          'password': {
            'from_secret': 'docker_password'
          }
        },
        'when': {
          'event': {
            'exclude': [
              'pull_request'
            ]
          }
        }
      }
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**'
      ]
    }
  }

def manifest():
  return [{
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'manifest',
    'steps': [
      {
        'name': 'manifest',
        'image': 'plugins/manifest',
        'pull': 'always',
        'settings': {
          'auto_tag': 'true',
          'username': {
            'from_secret': 'docker_username'
          },
          'password': {
            'from_secret': 'docker_password'
          },
          'spec': 'docker/manifest.tmpl',
          'ignore_missing': 'true',
        },
      },
      {
        'name': 'microbadger',
        'image': 'plugins/webhook',
        'pull': 'always',
        'settings': {
          'urls': {
            'from_secret': 'microbadger_url'
          }
        },
      }
    ],
    'depends_on': [],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**'
      ]
    }
  }]

def gitter():
  return [{
    'kind': 'pipeline',
    'type': 'docker',
    'name': 'gitter',
    'clone': {
      'disable': True
    },
    'steps': [
      {
        'name': 'gitter',
        'image': 'plugins/gitter',
        'pull': 'always',
        'settings': {
          'webhook': {
            'from_secret': 'gitter_webhook'
          }
        },
      },
    ],
    'depends_on': [
      'manifest'
    ],
    'trigger': {
      'ref': [
        'refs/heads/master',
        'refs/tags/**',
        'refs/pull/**'
      ],
      'status': [
        'failure'
      ]
    }
  }]
