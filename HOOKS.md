# Hooks

A hook is an executable file that the Addon-operator executes when some event occurs. It can be a script or a compiled program written in any programming language.

The Addon-operator pursues an agreement stating that the information is transferred to hooks via files and results of hooks execution are also stored in files. Paths to files are passed via environment variables. The output to stdout will be written to the log, except for the case with the configuration output (run with `--config` flag). Such an agreement simplifies the work with the input data and reporting the results of the hook execution.

## Global hook

Global hooks are stored in the $GLOBAL_HOOKS_DIR directory. The Addon-operator recursively searches all executable files in it and runs them with the `--config` flag. Each hook prints its events binding configuration in the JSON format to stdout. If the execution fails, the Addon-operator terminates with the code of 1.

Bindings from [shell-operator](https://github.com/flant/shell-operator) are available for global hooks: [onStartup](#onstartup), [schedule](#schedule) and [kubernetes](#kubernetes). The bindings to the events of the modules discovery process are also available: [beforeAll](#beforeall) and [afterAll](#afterall) (see [modules discovery](LIFECYCLE.md#modules-discovery)).

During execution, a global hook receives global values. These values can be modified by the hook to share data with global hooks, module hooks and helm templates. If hook changes global values, the 'global values changed' event is generated and then all modules are reloaded. For details on values storage, see [VALUES](VALUES.md). See also [an overview](LIFECYCLE.md#reload-all-modules) and [a detailed description](LIFECYCLE-STEPS.md#reload-all-modules) of 'Reload all modules' process.

## Module hook

Module hooks are executable files stored in the `hooks` subdirectory of the module. During the ['modules discovery'](LIFECYCLE.md#modules-discovery) process, if module appears to be enabled, the Addon-operator searches for executable files in `hooks` directory and all found files are executed with `--config` flag. Each hook prints its event binding configuration in JSON format to stdout. The module discovery process restarts if an error occurs.

Bindings from [shell-operator](https://github.com/flant/shell-operator) are available for module hooks: [schedule](#schedule) and [kubernetes](#kubernetes). The bindings of the module lifecycle are also available: `onStartup`, `beforeHelm`, `afterHelm`, `afterDeleteHelm` — see [module lifecycle](LIFECYCLE.md#module-lifecycle).

During execution, a module hook receives global values and module values. Module values can be modified by the hook to share data with other hooks of the same module. If hook changes module values, the 'module values changed' event is generated and then the module is reloaded. For details on values storage, see [VALUES](VALUES.md). See also a [module lifecycle](LIFECYCLE.md#module-lifecycle) and a [module run](LIFECYCLE-STEPS.md#module-run) detailed description.


# Bindings

## Overview

| Binding  | Global? | Module? | Info |
| ------------- | ------------- | --- | --- |
| [onStartup](#onstartup)↗  | ✓ | – | On Addon-operator startup |
| [onStartup](#onstartup)↗  | – |  ✓ | On first module run |
| [beforeAll](#beforeall)↗ | ✓ | – | Before run all modules |
| [afterAll](#afterall)↗ | ✓ | – | After run all modules |
| [beforeHelm](#beforehelm)↗ | – | ✓ | Before run helm install |
| [afterHelm](#afterhelm)↗ | – | ✓ | After helm install |
| [afterDeleteHelm](#afterdeletehelm)↗ | – | ✓ | After run helm delete |
| [schedule](#schedule)↗ | ✓ | ✓ | Run on schedule |
| [kubernetes](#kubernetes)↗ | ✓ | ✓ | Run on event from Kubernetes |

## onStartup

- Global hook execution on Addon-operator's start-up.
- Module hook execution on the first run of enabled module.

Syntax:

```
{
  "onStartup": ORDER
}
```

Parameters:
- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.

## beforeAll

The execution of global hooks before modules discovery.

Syntax:

```
{
  "beforeAll": ORDER
}
```

Parameters:

- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.

## afterAll

The execution of global hooks after running and removing modules.

Syntax:

```
{
  "afterAll": ORDER
}
```

Parameters:

- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.

## beforeHelm

The execution of a module hook before the helm chart installation (see [module lifecycle](LIFECYCLE.md#module-lifecycle)).

```
{
  "beforeHelm": ORDER
}
```

Parameters:
- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.

## afterHelm

The execution of a module hook after the helm chart installation (see [module lifecycle](LIFECYCLE.md#module-lifecycle)).

```
{
  "afterHelm": ORDER
}
```

Parameters:
- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.

## afterDeleteHelm

The execution of a module hook after the helm chart deletion (see [module lifecycle](LIFECYCLE.md#module-lifecycle)).

```
{
  "afterDeleteHelm": ORDER
}
```

Parameters:
- `ORDER` — the execution order (when added to the queue, the hooks will be sorted in the specified order, and then alphabetically). The value should be an integer.


## schedule

See the [schedule binding](https://github.com/flant/shell-operator/blob/master/HOOKS.md#schedule) from the Shell-operator.

## kubernetes

See the [kubernetes binding](https://github.com/flant/shell-operator/blob/master/HOOKS.md#kubernetes) from the Shell-operator.

> Note: Addon-operator requires a ServiceAccount with the appropriate [RBAC](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) permissions. See `addon-operator-rbac.yaml` files in [examples](/examples).

# Execution on event

When an event associated with a hook is triggered, Addon-operator executes the hook without arguments and passes the global or module values from the storage of the values via temporary files. In response, a hook could return JSON patches to modify values. The detailed description of the storage of the values is available in [VALUES](VALUES.md) document.

## Binding context

The binding context is an information about the event which caused the hook execution.

The $BINDING_CONTEXT_PATH environment variable contains the path to a file with a JSON array of structures with the following fields:

- `binding` is a string from the `name` parameter for `schedule` or `kubernetes` bindings. Its value is a *binding type* if the parameter is not set and for other hooks. For example, the binding context for `beforeAll` binding type:

```json
[{"binding":"beforeAll"}]
```

The binding context for `schedule` and `kubernetes` hooks contains additional fields, described in Shell-operator [documentation](https://github.com/flant/shell-operator/blob/master/HOOKS.md#binding-context).

`beforeAll` and `afterAll` global hooks and `beforeHelm`, `afterHelm`, and `afterDeleteHelm` module hooks are executed with the binding context that includes a `snapshots` field, which contains all Kubernetes objects that match hook's `kubernetes` bindings configurations. 

For example, a global hook with `kubernetes` and `beforeAll` bindings may have this configuration:

```
configVersion: v1
beforeAll: 10
kubernetes:
- name: monitor-pods
  apiVersion: v1
  kind: Pod
  jqFilter: ".metadata.labels"
```

This hook will be executed *before* updating the Helm release with this binding context:

```
[{"binding":"beforeAll",
"snapshots":{
  "monitor-pods":[
    {
      "object":{
        "kind":"Pod,
        "apiVersion": "v1",
        "metadata":{ "name":"pod-1r62e3", "namespace":"default", ...},
        ...
      },
      "filterResult": {
        "label1": "label value",
        ...
      },
    },
    ...
    more pods
    ...
  ]
}
}]
```
