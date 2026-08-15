package main

/*func Benchmark(b *testing.B) {
	envVars, err := loadEnvVars()
	if err != nil {
		b.Errorf("failed to load env vars: %v", err)
	}
	tok := envVars.discordToken
	client, err := disgo.New(tok,
		bot.WithDefaultGateway(),
	)
	if err != nil {
		b.Errorf("error while building disgo: %v", err)
	}
	defer client.Close(context.TODO())

	for b.Loop() {
	}
}*/
