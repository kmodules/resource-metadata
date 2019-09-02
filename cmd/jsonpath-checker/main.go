package main

import (
	"bytes"
	"fmt"
	"os"

	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/yaml"
)

/*
$ jsonpath-checker "{range .items[*]}{.metadata.name},{.status.capacity}{\"\n\"}{end}"
{range .items[*]}{.metadata.name},{.status.capacity}{"\n"}{end}
127.0.0.1,map[cpu:4]
127.0.0.2,map[cpu:8]

$ jsonpath-checker "{range .items[*]}{.metadata.name},{.status.capacity.cpu}{\"\n\"}{end}"
{range .items[*]}{.metadata.name},{.status.capacity.cpu}{"\n"}{end}
127.0.0.1,4
127.0.0.2,8
*/
var in = `{
  "kind": "List",
  "items":[
    {
      "kind":"None",
      "metadata":{"name":"127.0.0.1"},
      "status":{
        "capacity":{"cpu":"4"},
        "addresses":[{"type": "LegacyHostIP", "address":"127.0.0.1"}]
      }
    },
    {
      "kind":"None",
      "metadata":{"name":"127.0.0.2"},
      "status":{
        "capacity":{"cpu":"8"},
        "addresses":[
          {"type": "LegacyHostIP", "address":"127.0.0.2"},
          {"type": "another", "address":"127.0.0.3"}
        ]
      }
    }
  ],
  "users":[
    {
      "name": "myself",
      "user": {}
    },
    {
      "name": "e2e",
      "user": {"username": "admin", "password": "secret"}
    }
  ]
}`

/*
$ jsonpath-checker "{range .spec.volumes[*]}{.persistentVolumeClaim.claimName}{\"\n\"}{end}"
{range .spec.volumes[*]}{.persistentVolumeClaim.claimName}{"\n"}{end}
myclaim

*/
var yamlPod = `apiVersion: v1
kind: Pod
metadata:
  name: pv-recycler
  namespace: default
spec:
  restartPolicy: Never
  volumes:
  - name: mypd
    persistentVolumeClaim:
      claimName: myclaim
  - name: vol
    hostPath:
      path: /any/path/it/will/be/replaced
  containers:
  - name: pv-recycler
    image: "k8s.gcr.io/busybox"
    command: ["/bin/sh", "-c", "test -e /scrub && rm -rf /scrub/..?* /scrub/.[!.]* /scrub/*  && test -z \"$(ls -A /scrub)\" || exit 1"]
    volumeMounts:
    - name: vol
      mountPath: /scrub
`

// sa -> secrets
/*
$ jsonpath-checker "{range .secrets[*]}{.name}{\"\n\"}{end}"
{range .secrets[*]}{.name}{"\n"}{end}
default-token-d99rd
*/
var yamlServiceAccount = `apiVersion: v1
kind: ServiceAccount
metadata:
  creationTimestamp: "2019-08-11T01:06:46Z"
  name: default
  namespace: default
  resourceVersion: "361"
  selfLink: /api/v1/namespaces/default/serviceaccounts/default
  uid: 23f0716f-0552-46e5-9563-10bc3c1043c2
secrets:
- name: default-token-d99rd
`

// secrets -> s/a
// ref: https://github.com/kubernetes/client-go/blob/62f300f03a5f016c8befeae0abeb8faebd1a2fd2/util/jsonpath/jsonpath_test.go#L276
/*
$ jsonpath-checker "{.metadata.annotations.kubernetes\.io/service-account\.name}"
{.metadata.annotations.kubernetes\.io/service-account\.name}
default⏎                                                                                                                                                                                                           ~/g/s/g/a/kproxy (rd) $
*/
var yamlSecret = `apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNU1EZ3hNVEF4TURVMU1Wb1hEVEk1TURnd09EQXhNRFUxTVZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTUh5CjRlaVozYllqdGJ3UmhDU3ZXbFJYOUN2allsUWNsTktHQ0FpY2hmSmtyRU1qaEdzWXdVK3dCRFJlbG9tN3hXNXoKVkxzalpXNmovb3V6SWUzOXYvdDRNL2E5VTB5eGV0eTRLeUxHZ21Eb3BkVStaOTVKakk3N3YveWN3Q3ZuTGxNdgpCeUZMWUQxbkVlY0tSUTVQMHdKOG5PMVNvNkZHcmpmTFk2bDRaeUgyakZ3emVjWVVPdGVpKzk0ZHc4cHpRdGhJCkpLQnppTGhXSjhWOVBOU0c2bDhFMmdOdldTcFdPV3VrMDhvY2ZYeFMzQ0QzUzVDcnl1NnV3RENhemlHcVdsRmsKS2gwZWtpK1BwMjVIaUNZZFFrU29xMitpOURqdndqbVpiUjFXY0dQM2hhTmtPeHBScWo4NGQzajh0U3B2bnAzVAp6Wm0zRmpSNGFtSnBiRTVxVktNQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFKVnV1aFhHWUl3U0RBaWJHWUZ6S0pIdkduRUwKbEhYcUwwdzI1cVFrR3p4NUk0RUYwbzcxTmNqODZQVzFHT1JadndSRWJ6UFQ4Wi9wcHhoQkF5QWluREpDZEs0Zgp2NCtOQzZ2OEpidVlNRzVVQmEzTjJjd3V0eTcwSEszVkJQbkRNYit6cDd1WHhsdnNqZHRZcno1MEZHeFZ2MnFSCk5COXZmMHcxK1h3SEovZ2dTcWpqZkFvZno4SnBzTEU2ZDF6NndlMlNnZUxpSkxPODgrckYwNnNzdGtSemhiZFIKUTdKaENBTzd0R3dyRW1sUnVXMWRFMUZFb0poMVE2M010eGJkTXdHd3BONXY1VklOYzVZUUs1dE1VZGZKZHpxSwphM3ZKc0UvZmJIUjJUN2hpQS93QUJzalg4cmFLTlQwdUJ0VStvYUNLUHp5M0RVaWZvaTZsY1dxY2dIVT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  namespace: ZGVmYXVsdA==
  token: ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSmtaV1poZFd4MElpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WldOeVpYUXVibUZ0WlNJNkltUmxabUYxYkhRdGRHOXJaVzR0WkRrNWNtUWlMQ0pyZFdKbGNtNWxkR1Z6TG1sdkwzTmxjblpwWTJWaFkyTnZkVzUwTDNObGNuWnBZMlV0WVdOamIzVnVkQzV1WVcxbElqb2laR1ZtWVhWc2RDSXNJbXQxWW1WeWJtVjBaWE11YVc4dmMyVnlkbWxqWldGalkyOTFiblF2YzJWeWRtbGpaUzFoWTJOdmRXNTBMblZwWkNJNklqSXpaakEzTVRabUxUQTFOVEl0TkRabE5TMDVOVFl6TFRFd1ltTXpZekV3TkROak1pSXNJbk4xWWlJNkluTjVjM1JsYlRwelpYSjJhV05sWVdOamIzVnVkRHBrWldaaGRXeDBPbVJsWm1GMWJIUWlmUS5kOWo2Q1pQdFNmX0RqRTEzdzJxU25TMDl0Z1lBTGdlZ1hEREZIVElmek5iWng5OXo2cXlXS2pqVjhUaUtTOXl2b1hfbXZjMWVzRHdJbDNkbVJyalpHNmJYSkZlZUlZaktqTjJORjhTNmNMTGIxc29RRWFwRmxYT2xKLVptQ3RfdE5ITUNnckpxTFFOVjJ1WGFRVW43UG53VTY0ejBnWjRzWUxTZmFpRTNPeFlfTU5lXzBiZnhSa3NaM1RZTE9vVUE1Ty1oMjZ6NmFXT2wzYjNYTll0cXBoUFlTaUdKbndseDQtTmJyNG1BOU9qTjZFc0k0YUtOa0hJWUtxNGFpeDRZQzR1dk4yV2FSVklDU3hoS0czd0s2VGtvY3V1RnFlVFZPZW5MNVNRbzZsaGM3WURYRF9nSDdoVVM0Z29DdUNlS2ZuRmMtaXVfcU00cUg3SVhvdEh2Qmc=
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: default
    kubernetes.io/service-account.uid: 23f0716f-0552-46e5-9563-10bc3c1043c2
  creationTimestamp: "2019-08-11T01:06:46Z"
  name: default-token-d99rd
  namespace: default
  resourceVersion: "358"
  selfLink: /api/v1/namespaces/default/secrets/default-token-d99rd
  uid: e2561a0c-66f2-4aac-aff6-a63f7ca737d0
type: kubernetes.io/service-account-token
`

/*
n1
n2

n1,ns1
n2,ns2

n1,ns1,k1
n2,ns2,k2

n1,ns1,k1,apiGroup1
n2,ns2,k2,apiGroup2
*/

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: jsonpath-checker <path>")
		os.Exit(1)
	}

	var input interface{}
	err := yaml.Unmarshal([]byte(yamlSecret), &input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	name := "rd"
	allowMissingKeys := true
	template := os.Args[1]

	fmt.Println(template)

	j := jsonpath.New(name)
	j.AllowMissingKeys(allowMissingKeys)
	err = j.Parse(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "in %s, parse %s error %v", name, template, err)
		os.Exit(1)
	}
	buf := new(bytes.Buffer)
	err = j.Execute(buf, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "in %s, execute error %v", name, err)
		os.Exit(1)
	}
	fmt.Print(buf.String())
}
