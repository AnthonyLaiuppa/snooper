FROM scratch 
ADD main /
ADD sns.ini /
ADD ca-certificates.crt /etc/ssl/certs/
CMD ["/main"]
