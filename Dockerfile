FROM scratch
COPY ./website /website
EXPOSE 8080
CMD ["/website"]
