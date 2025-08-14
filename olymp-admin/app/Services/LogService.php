<?php

namespace App\Services;

class LogService
{
    private ElasticsearchService $elasticsearchService;
    public function __construct(
        ElasticsearchService $elasticsearchService
    )
    {
        $this->elasticsearchService = $elasticsearchService;
    }

    public function create($ip, $client, $url){$this->elasticsearchService->index('log', [
        'ip' => $ip,
        'url' => $url,
        'client' => $client,
        'timestamps' => date('Y-m-d H:i:s')
    ]);

        dd($this->elasticsearchService->decode($this->elasticsearchService->search('log', [
            'query' => [
                'match_all' => new \stdClass()  // Elasticsearch требует объект, а не массив
            ]
        ])));

    }
}
