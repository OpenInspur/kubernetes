FROM registry.icp.com:5000/library/os/neokylin-7.4-mips64el:5.0.0
ADD _output/bin/hyperkube /hyperkube
RUN ln -s /hyperkube /apiserver \
 && ln -s /hyperkube /cloud-controller-manager \
 && ln -s /hyperkube /controller-manager \
 && ln -s /hyperkube /kubectl \
 && ln -s /hyperkube /kubelet \
 && ln -s /hyperkube /proxy \
 && ln -s /hyperkube /scheduler \
 && ln -s /hyperkube /usr/local/bin/cloud-controller-manager \
 && ln -s /hyperkube /usr/local/bin/kube-apiserver \
 && ln -s /hyperkube /usr/local/bin/kube-controller-manager \
 && ln -s /hyperkube /usr/local/bin/kube-proxy \
 && ln -s /hyperkube /usr/local/bin/kube-scheduler \
 && ln -s /hyperkube /usr/local/bin/kubectl \
 && ln -s /hyperkube /usr/local/bin/kubelet
