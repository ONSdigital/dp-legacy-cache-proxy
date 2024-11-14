job "dp-legacy-cache-proxy" {
  datacenters = ["eu-west-2"]
  region      = "eu"
  type        = "service"

  update {
    stagger          = "60s"
    min_healthy_time = "30s"
    healthy_deadline = "2m"
    max_parallel     = 1
    auto_revert      = true
  }

  group "web" {
    count = "{{WEB_TASK_COUNT}}"

    spread {
      attribute = "${node.unique.id}"
      weight    = 100
      # with `target` omitted, Nomad will spread allocations evenly across all values of the attribute.
    }
    spread {
      attribute = "${attr.platform.aws.placement.availability-zone}"
      weight    = 100
      # with `target` omitted, Nomad will spread allocations evenly across all values of the attribute.
    }
    
    constraint {
      attribute = "${node.class}"
      value     = "web"
    }

    restart {
      attempts = 3
      delay    = "15s"
      interval = "1m"
      mode     = "delay"
    }

    network {
      port "http" {
        to = 29200
      }
    }

    service {
      name = "dp-legacy-cache-proxy"
      port = "http"
      tags = ["web"]

      check {
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "dp-legacy-cache-proxy-web" {
      driver = "docker"

      config {
        command = "./dp-legacy-cache-proxy"
        image   = "{{ECR_URL}}:concourse-{{REVISION}}"
        ports   = ["http"]
      }

      resources {
        cpu    = "{{WEB_RESOURCE_CPU}}"
        memory = "{{WEB_RESOURCE_MEM}}"
      }

      template {
        data = <<EOH
        # Configs based on environment (e.g. export BIND_ADDR=":{{ env "NOMAD_PORT_http" }}")
        # or static (e.g. export BIND_ADDR=":8080")

        export BABBAGE_URL="http://{{ env "NOMAD_IP_http" }}:10000"
        export LEGACY_CACHE_API_URL="http://{{ env "NOMAD_IP_http" }}:13850"
        export RELEASE_CALENDAR_URL="http://{{ env "NOMAD_IP_http" }}:13000"
        export SEARCH_CONTROLLER_URL="http://{{ env "NOMAD_IP_http" }}:11750"

        # Secret configs read from vault
        {{ with (secret (print "secret/" (env "NOMAD_TASK_NAME"))) }}
        {{ range $key, $value := .Data }}
        export {{ $key }}={{ $value | toJSON }}
        {{ end }}
        {{ end }}
        EOH

        destination = "secrets/app.env"
        env         = true
        splay       = "1m"
        change_mode = "restart"
      }

      vault {
        policies = ["dp-legacy-cache-proxy-web"]
      }
    }
  }
}
