encapsulating for the basic [k8s api](http://kubernetes.io/third_party/swagger-ui/#!/v1beta1/listNode)

using func Sendapi to send the k8s api

examples:

    using the API  GET /api/v1beta1/minions 
    getcommands := []string{"minions"}
    status, message := sendapi.Sendapi("GET", masterip, masterport, "v1beta1", getcommands, "")
	
    using the API  POST /api/v1beta1/replicationControllers
    patchcommands := []string{"replicationControllers", "frontend-controller"}
    status, message := sendapi.Sendapi("GET", masterip, masterport, "v1beta1", getcommands, filename)

the type of message is map[string]string

