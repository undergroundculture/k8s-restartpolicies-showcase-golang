# k8s-restartpolicies-showcase-go-crud

I've created a simple app that will be deployed to k8s and will have liveness, readiness and startup probes.

In the `/app` folder you can find the source code of the app with Dockerfile. There are timer that runs every 15seconds and randomize numbers:
1) if number is 10 then app will exit with code 1
2) if number is 1 then app will exit with code 0
3) if number not 10 or 1 then app will continue to run

> [!WARNING]
> *Application quite bloated since i was thinking wrong about restart policies - it's pod level mechanis primarally and apply it why running you application as deployment and quite useless (here i'm talking only about regular container inside pod without init and sidecar containers).*
> *I left code such as to remind myself where i was wrong. *

In the `/manifests` folder you can find the deployment manifest with commented restartPolcies. You can make sure that they are not working no matter how container exist - replica set will make sure that there are desired number of containers are running.
There are also /pod manifset to practice with.

The image is `undergroundculture/restart-policies:1.0` and it can be easilly pulled from docker hub.


## restartPolicy: Always

### zero exit code

App logs on exit with code 0:
```
~/: kubectl logs -f pod/restart-policies
2026/07/23 17:14:14.135454 PID: 1
2026/07/23 17:14:14.634687 doing initialization
2026/07/23 17:14:29.648305 initialization done
2026/07/23 17:14:29.648525 HTTP server listening on :8080
2026/07/23 17:14:39.657741 timer: generated number 1
2026/07/23 17:14:39.657780 random 1 -> exiting with code 0
```


```bash
~: kubectl get pods -w -l app=restart-policies
NAME               READY   STATUS    RESTARTS   AGE
restart-policies   0/1     Pending   0          0s
restart-policies   0/1     Pending   0          0s
restart-policies   0/1     ContainerCreating   0          0s
restart-policies   0/1     ContainerCreating   0          6s
restart-policies   1/1     Running             0          11s
restart-policies   0/1     Completed           0          38s
restart-policies   1/1     Running             1 (3s ago)   40s
```


### non-zero exit code

```bash
2026/07/23 17:17:41.343708 timer: generated number 9
2026/07/23 17:17:51.343822 timer: generated number 2
2026/07/23 17:18:01.343930 timer: generated number 7
2026/07/23 17:18:11.344040 timer: generated number 4
2026/07/23 17:18:21.344118 timer: generated number 9
2026/07/23 17:18:31.344208 timer: generated number 3
2026/07/23 17:18:41.336492 timer: generated number 1
2026/07/23 17:18:41.336531 random 10 -> exiting with code 1

``` 


```bash
~: kubectl get pods -w -l app=restart-policies
NAME               READY   STATUS              RESTARTS   AGE
restart-policies   0/1     ContainerCreating   0          3s
restart-policies   0/1     ContainerCreating   0          8s
restart-policies   1/1     Running             0          15s
restart-policies   1/1     Running             0          93s
restart-policies   0/1     Completed           0          3m2s
restart-policies   1/1     Running             1 (6s ago)   3m7s
restart-policies   0/1     Completed           1 (43s ago)   3m44s
restart-policies   0/1     CrashLoopBackOff    1 (12s ago)   3m55s
restart-policies   1/1     Running             2 (15s ago)   3m58s
restart-policies   1/1     Terminating         2 (27s ago)   4m10s
restart-policies   1/1     Terminating         2 (27s ago)   4m10s

```


## restartPolicy: OnFailure

### zero exit code 

```bash
~/: sleep 15 && kubectl logs -f pod/restart-policies
2026/07/23 17:22:42.588475 timer: generated number 4
2026/07/23 17:22:52.593663 timer: generated number 8
2026/07/23 17:23:02.592518 timer: generated number 4
2026/07/23 17:23:12.584525 timer: generated number 1
2026/07/23 17:23:12.584564 random 1 -> exiting with code 0
~/: kubectl get pods
NAME                                  READY   STATUS      RESTARTS      AGE
go-crud-with-probes-8b94f5494-2fvpn   1/1     Running     1 (20h ago)   24h
go-crud-with-probes-8b94f5494-cfck2   1/1     Running     1 (20h ago)   24h
go-crud-with-probes-8b94f5494-z85w9   1/1     Running     1 (20h ago)   24h
grafana-dc68c9898-8xkbh               1/1     Running     2 (20h ago)   4d4h
restart-policies                      0/1     Completed   0             83s
```

```bash
~: kubectl get pods -w -l app=restart-policies
NAME               READY   STATUS    RESTARTS   AGE
restart-policies   0/1     Pending   0          0s
restart-policies   0/1     Pending   0          0s
restart-policies   0/1     ContainerCreating   0          0s
restart-policies   0/1     ContainerCreating   0          5s
restart-policies   1/1     Running             0          10s
restart-policies   0/1     Completed           0          66s
restart-policies   0/1     Completed           0          68s
```
 

### non-zero exit 

```bash
~/: kubectl logs -f pod/restart-policies
2026/07/23 17:24:52.159945 PID: 1
2026/07/23 17:24:52.460159 doing initialization
2026/07/23 17:25:07.464357 initialization done
2026/07/23 17:25:07.464871 HTTP server listening on :8080
2026/07/23 17:25:07.559972 received signal: urgent I/O condition
2026/07/23 17:25:07.560102 no special handling for urgent I/O condition
2026/07/23 17:25:07.659553 received signal: urgent I/O condition
2026/07/23 17:25:07.659779 no special handling for urgent I/O condition
2026/07/23 17:25:17.474008 timer: generated number 6
2026/07/23 17:25:27.465043 timer: generated number 2
2026/07/23 17:25:37.468373 timer: generated number 10
2026/07/23 17:25:37.468430 random 10 -> exiting with code 1
~/: kubectl get pods
NAME                                  READY   STATUS    RESTARTS      AGE
go-crud-with-probes-8b94f5494-2fvpn   1/1     Running   1 (20h ago)   24h
go-crud-with-probes-8b94f5494-cfck2   1/1     Running   1 (20h ago)   24h
go-crud-with-probes-8b94f5494-z85w9   1/1     Running   1 (20h ago)   24h
grafana-dc68c9898-8xkbh               1/1     Running   2 (20h ago)   4d4h
restart-policies                      1/1     Running   1 (34s ago)   92s

```

```bash
~: kubectl get pods -w -l app=restart-policies
NAME               READY   STATUS              RESTARTS   AGE
restart-policies   0/1     ContainerCreating   0          2s
restart-policies   0/1     ContainerCreating   0          6s
restart-policies   1/1     Running             0          12s
restart-policies   0/1     Error               0          58s
restart-policies   1/1     Running             1 (4s ago)   62s
```


## RestartPolicy: Never

App logs:

```bash
~/: kubectl logs -f pod/restart-policies
2026/07/23 17:09:53.095622 PID: 1
2026/07/23 17:09:54.201067 doing initialization
2026/07/23 17:10:09.212449 initialization done
2026/07/23 17:10:09.213116 HTTP server listening on :8080
2026/07/23 17:10:19.222698 timer: generated number 4
2026/07/23 17:10:29.222789 timer: generated number 8
2026/07/23 17:10:39.222904 timer: generated number 8
2026/07/23 17:10:49.222981 timer: generated number 10
2026/07/23 17:10:49.223017 random 10 -> exiting with code 1
```


```bash
~: kubectl get pods -w -l app=restart-policies
NAME                               READY   STATUS        RESTARTS   AGE
restart-policies-d575fd464-grcfd   1/1     Terminating   0          52s
restart-policies-d575fd464-grcfd   0/1     Error         0          55s
restart-policies-d575fd464-grcfd   0/1     Error         0          56s
restart-policies-d575fd464-grcfd   0/1     Error         0          56s
restart-policies-d575fd464-grcfd   0/1     Error         0          56s
restart-policies                   0/1     Pending       0          0s
restart-policies                   0/1     Pending       0          0s
restart-policies                   0/1     ContainerCreating   0          0s
restart-policies                   0/1     ContainerCreating   0          9s
restart-policies                   1/1     Running             0          12s
restart-policies                   0/1     Error               0          69s
restart-policies                   0/1     Error               0          71s
restart-policies                   0/1     Error               0          71s

```
