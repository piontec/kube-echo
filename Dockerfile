FROM scratch
COPY --chmod=755 kube-echo /kube-echo
USER 65534:65534
ENTRYPOINT [ "/kube-echo" ]
