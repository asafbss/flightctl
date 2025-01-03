# Auto-Registering Devices with MicroShift into ACM

## Pre-requisite

The Operating System image running on our fleet of devices must contain MicroShift. This is essential for running the Klusterlet, the agent responsible for communication with Red Hat Advanced Cluster Management (ACM)

## Configuring ACM Agent Registration

Auto-registration of devices with MicroShift into ACM relies on a feature called agent registration, allowing clusters to be imported via REST API calls using a CA bundle and token. Follow these [instructions](https://docs.redhat.com/en/documentation/red_hat_advanced_cluster_management_for_kubernetes/2.11/html/clusters/cluster_mce_overview#importing-managed-agent) to configure the necessary RBAC policies for your own user. Skip running the command to import the managed clusters, as this step will be automatically handled by the fleet devices.

As part of this setup, you will need to collect the following information:

- Agent registration URL
- ca.crt
- Token

## Creating Repositories in Flight Control

A new repository must be created with the data collected before from ACM's agent registration endpoint. This repository will be used by the fleet configuration using the HTTP config provider described in the following [link](https://github.com/flightctl/flightctl/blob/main/docs/user/managing-devices.md#getting-configuration-from-an-http-server).

A sample YAML file for an HTTP repository looks like:

```console
apiVersion: v1alpha1
kind: Repository
metadata:
  labels: {}
  name: acm-registration
spec:
  httpConfig:
    token: $token
    ca.crt: <base64-encoded ca.crt>
  type: http
  url: https://$agent_registration_host
  validationSuffix: /agent-registration/crds/v1
```

## Fleet Definition Overview

The following fleet definition provides an example of how devices running MicroShift can be auto-registered into an ACM hub. Here is a breakdown of its sections:

```console
apiVersion: v1alpha1
kind: Fleet
metadata:
  labels: {}
  name: fleet-acm
spec:
  selector:
    matchLabels:
      fleet: acm
  template:
    metadata:
      generation: 1
      labels:
        fleet: fleet-acm
    spec:
      os:
        image: quay.io/myorg/device-image-with-microshift:v1
      config:
      - name: acm-crd
        httpRef:
          filePath: /var/local/acm-import/crd.yaml
          repository: acm-registration
          suffix: /agent-registration/crds/v1
      - name: acm-import
        httpRef:
          filePath: /var/local/acm-import/import.yaml
          repository: acm-registration
          suffix: /agent-registration/manifests/{{ device.metadata.name }}
      - name: pull-secret
        inline:
            - path: "/etc/crio/openshift-pull-secret"
              content: "{\"auths\":{...}}"
      hooks:
        afterUpdating:
          - path: "/var/local/acm-import/crd.yaml"
            onFile: [Create]
            actions:
              - executable:
                  run: "kubectl apply -f /var/local/acm-import/crd.yaml"
                  envVars: ["KUBECONFIG=/var/lib/microshift/resources/kubeadmin/kubeconfig"]
          - path: "/var/local/acm-import/import.yaml"
            onFile: [Create]
            actions:
              - executable:
                  run: "kubectl apply -f /var/local/acm-import/import.yaml"
                  envVars: ["KUBECONFIG=/var/lib/microshift/resources/kubeadmin/kubeconfig"]
```

### Fleet Specification Breakdown

As described in the user documentation, a fleet specification is composed of various sections. Let us deep dive in the config and hooks section of our sample fleet.

1. **Configuration**: The configuration section uses the HTTP configuration provider to fetch information from an endpoint. The repository `acm-registration` contains the registration URL for ACM's agent registration:

    ```console
    - name: acm-crd
      httpRef:
        filePath: /var/local/acm-import/crd.yaml
        repository: acm-registration
        suffix: /agent-registration/crds/v1
    ```

    This retrieves the CRD from the ACM endpoint `https://$agent_registration_host/agent-registration/crds/v1` and stores it at the specified file path.

    The next configuration retrieves the cluster import manifests. As shown below, the HttpConfigProviderSpec supports Flight Control template mechanism, so the device name can be used as part of the suffix.

    ```console
    - name: acm-import
      httpRef:
        filePath: /var/local/acm-import/import.yaml
        repository: acm-registration
        suffix: /agent-registration/manifests/{{ device.metadata.name }}
    ```

    This API call to ACM's agent registration endpoint will retrieve a set of Kubernetes manifests to deploy the Klusterlet. Once we have both the Klusterlet CRD and deployment manifests, we will use Flight Control hooks to apply them.

    ```console
    - name: pull-secret
      inline:
          - path: "/etc/crio/openshift-pull-secret"
            content: "{\"auths\":{...}}"
    ```

    This adds your OpenShift pull secret to the device.  You can download your pull secret [here](https://cloud.redhat.com/openshift/install/pull-secret).  If you do not wish to have your pull secret available via the Flight Control API, you may store it in a private Git repository or a Kubernetes Secret, and change the fleet specification accordingly.

1. **Hooks**: Once the configuration is fetched, the hooks apply the Kubernetes manifests to the device:

```console
hooks:
  afterUpdating:
    - path: "/var/local/acm-import/crd.yaml"
      onFile: [Create]
      actions:
        - executable:
            run: "kubectl apply -f /var/local/acm-import/crd.yaml"
            envVars: ["KUBECONFIG=/var/lib/microshift/resources/kubeadmin/kubeconfig"]
            workDir: "/usr/bin/"
    - path: "/var/local/import.yaml"
      onFile: [Create]
      actions:
        - executable:
            run: "kubectl apply -f /var/local/acm-import/import.yaml"
            envVars: ["KUBECONFIG=/var/lib/microshift/resources/kubeadmin/kubeconfig"]
            workDir: "/usr/bin/"
```

Flight Control hooks are explained in the following [link](https://github.com/flightctl/flightctl/blob/main/docs/user/managing-devices.md#using-device-lifecycle-hooks).
Hooks are triggered at specific moments on the lifecycle of a device such as before/after updating, before/after reboot, etc.
The hooks definition shown above describe two hooks executed after updating the fleet, one to apply the CRD file and one to apply the rest of the Klusterlet manifests stored in the import.yaml file.

## Device Reconciliation Process

Once devices are enrolled in Flight Control and assigned to a fleet definition like the one above, the agent will begin a reconciliation process as outlined in the project documentation. This process ensures each device is updated with the required OS image, configuration and applications. This specific workflow  ensures all steps required to register a MicroShift cluster into ACM are performed. From updating the Operating System image to one containing MicroShift and its dependencies, getting the Klusterlet manifests into disk and apply them.

Once the manifests are applied by the hooks mechanism, MicroShift will run the Klusterlet agent which will contact the ACM hub and send a Certificate Signing Request, as part of its registration process.

This automated workflow allows for the large-scale, hands-free registration of thousands of MicroShift clusters into ACM.

## Summary

This guide outlines an automated process to register MicroShift clusters running on Flight Control managed devices into ACM at scale. By leveraging Flight Control’s configuration providers, hooks and templates, users can import MicroShift instances into ACM seamlessly and securely, streamlining cluster management across a fleet.
