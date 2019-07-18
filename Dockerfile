FROM scratch 
ADD snooper_unix /
ADD sns.ini /
ADD ca-certificates.crt /etc/ssl/certs/
CMD ["/snooper_unix"]
