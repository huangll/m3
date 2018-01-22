	"github.com/m3db/m3db/environment"
	"github.com/m3db/m3db/x/mmap"
	runtimeOpts := m3dbruntime.NewOptions().
		SetPersistRateLimitOptions(ratelimit.NewOptions().
			SetLimitEnabled(true).
			SetLimitMbps(cfg.Filesystem.ThroughputLimitMbps).
			SetLimitCheckEvery(cfg.Filesystem.ThroughputCheckEvery)).
		SetWriteNewSeriesAsync(cfg.WriteNewSeriesAsync).
		SetWriteNewSeriesBackoffDuration(cfg.WriteNewSeriesBackoffDuration)

	if tick := cfg.Tick; tick != nil {
		runtimeOpts = runtimeOpts.
			SetTickSeriesBatchSize(tick.SeriesBatchSize).
			SetTickPerSeriesSleepDuration(tick.PerSeriesSleepDuration).
			SetTickMinimumInterval(tick.MinimumInterval)
	}

	runtimeOptsMgr := m3dbruntime.NewOptionsManager()
	if err := runtimeOptsMgr.Update(runtimeOpts); err != nil {
		logger.Fatalf("could not set initial runtime options: %v", err)
	}
	defer runtimeOptsMgr.Close()

	opts = opts.SetRuntimeOptionsManager(runtimeOptsMgr)

		logger.Fatalf("could not parse new file mode: %v", err)
		logger.Fatalf("could not parse new directory mode: %v", err)
	}

	mmapCfg := cfg.Filesystem.MmapConfiguration()
	shouldUseHugeTLB := mmapCfg.HugeTLB.Enabled
	if shouldUseHugeTLB {
		// Make sure the host supports HugeTLB before proceeding with it to prevent
		// excessive log spam.
		shouldUseHugeTLB, err = hostSupportsHugeTLB()
		if err != nil {
			logger.Fatalf("could not determine if host supports HugeTLB: %v", err)
		}
		if !shouldUseHugeTLB {
			logger.Warnf("host doesn't support HugeTLB, proceeding without it")
		}
		SetSeekReaderBufferSize(cfg.Filesystem.SeekReadBufferSize).
		SetMmapEnableHugeTLB(shouldUseHugeTLB).
		SetMmapHugeTLBThreshold(mmapCfg.HugeTLB.Threshold).
		SetRuntimeOptionsManager(runtimeOptsMgr)
		logger.Fatalf("unknown commit log queue size type: %v",
	// Set the series cache policy
	seriesCachePolicy := cfg.Cache.SeriesConfiguration().Policy
	opts = opts.SetSeriesCachePolicy(seriesCachePolicy)
	// Apply pooling options
	// Setup the block retriever
	switch seriesCachePolicy {
	case series.CacheAll:
		// No options needed to be set
	default:
		// All other caching strategies require retrieving series from disk
		// to service a cache miss
		retrieverOpts := fs.NewBlockRetrieverOptions().
			SetBytesPool(opts.BytesPool()).
			SetSegmentReaderPool(opts.SegmentReaderPool()).
			SetIdentifierPool(opts.IdentifierPool())
		if blockRetrieveCfg := cfg.BlockRetrieve; blockRetrieveCfg != nil {
			retrieverOpts = retrieverOpts.
				SetFetchConcurrency(blockRetrieveCfg.FetchConcurrency)
		}
		blockRetrieverMgr := block.NewDatabaseBlockRetrieverManager(
			func(md namespace.Metadata) (block.DatabaseBlockRetriever, error) {
				retriever := fs.NewBlockRetriever(retrieverOpts, fsopts)
				if err := retriever.Open(md); err != nil {
					return nil, err
				}
				return retriever, nil
			})
		opts = opts.SetDatabaseBlockRetrieverManager(blockRetrieverMgr)
	}
	pm, err := fs.NewPersistManager(fsopts)
		logger.Fatalf("could not create persist manager: %v", err)
	opts = opts.SetPersistManager(pm)
	var (
		envCfg environment.ConfigureResults
	)
	switch {
	case cfg.EnvironmentConfig.Service != nil:
		logger.Info("creating dynamic config service client with m3cluster")

		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{
			InstrumentOpts: iopts,
			HashingSeed:    cfg.Hashing.Seed,
		})
		if err != nil {
			logger.Fatalf("could not initialize dynamic config: %v", err)
		}

	case cfg.EnvironmentConfig.Static != nil:
		logger.Info("creating static config service client with m3cluster")

		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{})
		if err != nil {
			logger.Fatalf("could not initialize static config: %v", err)
		}

	default:
		logger.Fatal("config service or static configuration required")
	}

	opts = opts.SetNamespaceInitializer(envCfg.NamespaceInitializer)

	topo, err := envCfg.TopologyInitializer.Init()
			TopologyInitializer: envCfg.TopologyInitializer,
	bs, err := cfg.Bootstrap.New(opts, m3dbClient)
	kvWatchBootstrappers(envCfg.KVStore, logger, timeout, cfg.Bootstrap.Bootstrappers,
			updated, err := cfg.Bootstrap.New(opts, m3dbClient)
	db, err := cluster.NewDatabase(hostID, envCfg.TopologyInitializer, opts)
		kvWatchNewSeriesLimitPerShard(envCfg.KVStore, logger, topo,
	err = setNewSeriesLimitPerShardOnChange(topo, runtimeOptsMgr, initClusterLimit)
	if err != nil {
		logger.Warnf("unable to set cluster new series insert limit: %v", err)
	}
			err = setNewSeriesLimitPerShardOnChange(topo, runtimeOptsMgr, value)
			if err != nil {
				logger.Warnf("unable to set cluster new series insert limit: %v", err)
				continue
			}
) error {
		return nil
	newRuntimeOpts := runtimeOpts.
		SetWriteNewSeriesLimitPerShardPerSecond(perPlacedShardLimit)
	return runtimeOptsMgr.Update(newRuntimeOpts)
	iopts := opts.InstrumentOptions()
	bytesPoolOpts := pool.NewObjectPoolOptions().
		SetInstrumentOptions(iopts.SetMetricsScope(scope.SubScope("bytes-pool")))
	checkedBytesPoolOpts := bytesPoolOpts.
		SetInstrumentOptions(iopts.SetMetricsScope(scope.SubScope("checked-bytes-pool")))
		b.Count = bucket.Size
		b.Options = bytesPoolOpts.
			SetRefillLowWatermark(bucket.RefillLowWaterMark).
			SetRefillHighWatermark(bucket.RefillHighWaterMark)
		logger.Infof("bytes pool registering bucket capacity=%d, size=%d, "+
			"refillLowWatermark=%f, refillHighWatermark=%f",
			bucket.Capacity, bucket.Size,
			bucket.RefillLowWaterMark, bucket.RefillHighWaterMark)
		bytesPool = pool.NewCheckedBytesPool(
			buckets,
			checkedBytesPoolOpts,
			func(s []pool.Bucket) pool.BytesPool {
				return pool.NewBytesPool(s, bytesPoolOpts)
			})
		bytesPool = pool.NewCheckedBytesPool(
			buckets,
			checkedBytesPoolOpts,
			func(s []pool.Bucket) pool.BytesPool {
				return pool.NewNativeHeap(s, bytesPoolOpts)
			})
	logger.Infof("bytes pool %s init", policy.Type)

func hostSupportsHugeTLB() (bool, error) {
	// Try and determine if the host supports HugeTLB in the first place
	withHugeTLB, err := mmap.Bytes(10, mmap.Options{
		HugeTLB: mmap.HugeTLBOptions{
			Enabled:   true,
			Threshold: 0,
		},
	})
	if err != nil {
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	}
	defer mmap.Munmap(withHugeTLB.Result)

	if withHugeTLB.Warning == nil {
		// If there was no warning, then the host didn't complain about
		// usa of huge TLB
		return true, nil
	}

	// If we got a warning, try mmap'ing without HugeTLB
	withoutHugeTLB, err := mmap.Bytes(10, mmap.Options{})
	if err != nil {
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	}
	defer mmap.Munmap(withoutHugeTLB.Result)
	if withoutHugeTLB.Warning == nil {
		// The machine doesn't support HugeTLB, proceed without it
		return false, nil
	}
	// The warning was probably caused by something else, proceed using HugeTLB
	return true, nil
}