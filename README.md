# cactus-maestro

Cactus On Kube - Maestro is a concept for turning decorated functions into Kubernetes-ready workloads. The goal is to ingest source files, validate function metadata, build an execution graph, package each function into its own container image, and expose selected functions through a scheduler-driven entrypoint.

## Flow

### 1. Ingest source code

The system starts by ingesting a file, such as a Python module.

```py
from cactuskit import ApiMethod, ApiProtocol, HttpStatus, cactuize

@cactuize()
def simple_entrypoint():
    return f"Hello World from {simple_entrypoint}"

@cactuize(
    auth=authenticate,
    protocol=ApiProtocol.HTTP,
    method=ApiMethod.GET,
)
def entrypoint(name):
    return HttpStatus.HTTP_CUSTOM(201), {
        "content": f"Hello {name}",
    }
```

### 2. Validate configuration

Validate each `cactuize` declaration before any further processing. If validation fails, stop immediately and return a structured error.

Example response:

```json
{
  "status": "error",
  "message": "ln. 5: error [XXX]: why it does not work"
}
```

### 3. Parse and build a graph

Parse each function, determine how functions relate to one another, and build a dependency graph that includes:

- direct dependencies
- helpers
- environment variables
- other runtime inputs

### 4. Encapsulate

Build a dedicated Docker image for each function. The image should be created from the graph produced by the parser and then pushed to the local Docker registry.

### 5. Expose as an entrypoint

Only functions decorated with `@cactuize` are exposed.

If `protocol` and `method` are not provided, the default configuration is:

```json
{
  "protocol": "HTTP",
  "protocol_cfg": {
    "method": "GET"
  }
}
```

If `protocol` and `method` are provided, those values are used as-is.

### 6. Schedule execution

Once a decorated function is registered as an entrypoint, it becomes publicly reachable. A request to that entrypoint publishes an event containing details such as the path, function name, and arguments.

Scheduler workers consume those events and handle execution:

1. Verify that the function exists.
2. Verify that the target node has enough resources.
3. Launch the correct image as a pod on the cluster.
4. Wait for the pod to finish.
5. Return the output as the response body.

## Intended outcome

The result is a function-to-cluster execution pipeline where source code is validated, packaged, scheduled, and served with minimal manual wiring.


[text](cmaestro-api) [text](cmaestro-db) [text](cmaestro-encapsulate) [text](cmaestro-entrypoint) [text](cmaestro-ingest) [text](cmaestro-parser) [text](cmaestro-scheduler) [text](cmaestro-validation)