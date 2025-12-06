RUN cd swisseph \
    && cc -O2 -fPIC -c sweph.c swedate.c swehouse.c swejpl.c swemmoon.c swemplan.c swephlib.c swecl.c swehel.c \
    && cc -shared -o libswe.so sweph.o swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o swephlib.o swecl.o swehel.o -lm \
    && cp libswe.so /usr/local/lib/ \
    && cp swephexp.h /usr/local/include/ \
    && ldconfig
