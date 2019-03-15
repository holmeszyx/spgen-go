gofile= spgen.go \
		genertor.go \
		androidgen.go \
		newer.go

spgen: $(gofile)
	go build -o spgen $(gofile)

run: spgen
	./spgen testdata/sp_config.toml

kt: spgen
	./spgen -o android:kt testdata/sp_config.toml

clean:
	-rm spgen
	-rm -rf testdata/abc