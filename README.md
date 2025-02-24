## NAME
doocot - ad hoc, securely.

## INSTALL

If you have Go installed:

```
go install github.com/thisdougb/doocot
```

Otherwise compiled releases for Mac, Linux, and FreeBSD are available at [doocot.sh/releases](https://doocot.sh/releases).

## SYNOPSIS

<pre>
Usage:
  doocot put [-v] [-once] [-words] [-json] { -create <i>n</i> | <i>data</i> }
  doocot get [-raw] <i>id</i>
</pre>

## DESCRIPTION

The **doocot** utility provides a convenient interface to the Doocot API. It simplifies securely sharing ad-hoc data between people and systems.

Using the free public backend, <i>data</i> is limited to 100 bytes.
This is enough for sharing API keys, passwords, etc.
It is encrypted in transit and at-rest whilst in the backend, then expired permanently after 15 minutes.

The **put** subcommand returns an identifier that is used to retrieve the data. This is the only instance of the (proxy) decryption key in existance. It is not possible to recover the encrypted data if the relevant identifier is lost.

To be CLI friendly, the response output is an easy copy/paste as a **doocot** **get** command and as a curl command.

The following options are available:

<pre>
  <strong>-v</strong>           verbose commentary

  <strong>-create <i>n</i></strong>    generate a new secret of length <i>n</i>, where 0 < <i>n</i> <= 100
                                                                         
  <strong>-once </n></strong>       expire the data immediately after it is first accessed

  <strong>-words</strong>       return a more human-friendly passphrase link rather than hex string

  <strong>-json</strong>        output in json format

  <strong>-raw</strong>         get the raw data from the backend, if your curious about secure storage
</pre>


Secrets in the public Doocot backend are expired automatically, after 15 minutes. 

There is no data retention. Using **doocot** to retrieve an expired (or non-existant) secret returns "Not Found".

### Doocot Backend Service

The free public backend service that **doocot** uses by default can be found [here](https://doocot.sh).

It is principally offered as a try *before you buy* service, but remains the full implementation of the doocot-server.
As a free public service though, it is strictly rate limited.
But we have tried to implement rate limiting in such a way as to be invisible to fair-use behaviour.

## EXIT STATUS

The doocot utility exits 0 on success, and > 0 if an error occurs.

## EXAMPLES

### One to another

Alice wants to send Bob an API key for a test account of some new service:

```
alice $ doocot put -words the api key is 564231f5bbe0a7e2833fe6dc1b66a40e9d7960229cd7040b23b4c2d4bf6eec43
slight-step-zoo-flock
```

Alice messages Bob the id 'slight-step-zoo-flock', using a chat service such as Slack or Teams.
Bob retrieves the api key:

```
bob $ doocot get slight-step-zoo-flock
the api key is 564231f5bbe0a7e2833fe6dc1b66a40e9d7960229cd7040b23b4c2d4bf6eec43
```

### Safer CI and automation

Creating system passwords should be really easy, but also easy to keep them out of your DevOps tool logs.

Here, task1 asks the backend to automatically generate a new password and parses the json response to get the identifier url.
It doesn't matter how high the logging is, this url is useless after the fact.

```
task1 $ new_password_url=$(doocot put -once -json -create 20 | jq -r '.url'); echo $new_password_url
https://doocot.sh/api/data/6e77fcd193db295a23254fedf41ee2c5dcf69eeb7401445d08f7d7d947d96419
```

Task2 can read the new password without worrying about leaking sensitive data to logs.
On creation the **-once** flag was set, so the secret is expired after it is read the first time.

```
task2 $ new_password=$(curl "$new_password_url"); echo $new_password
wSDvXq9AT5EIhoqelzg0
```

*(The echo command here is just for clarity.)*

## SEE ALSO

The public-use backend service at [doocot.sh](https://doocot.sh).

Compiled releases for Mac, Linux, and FreeBSD are availble at [doocot.sh/releases](https://doocot.sh/releases).

## AUTHORS

[thisdougb](https://github.com/thisdougb)

&copy; Far Oeuf Limited, 2025
