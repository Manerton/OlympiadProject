<?php

namespace App\Services;

class LogService
{
    public const LOG = 'log';
    private ElasticsearchService $elasticsearchService;
    public function __construct(
        ElasticsearchService $elasticsearchService
    )
    {
        $this->elasticsearchService = $elasticsearchService;
    }

    public function create($ip, $client, $url)
    {
        $params = [
            'index' => self::LOG,
            'body' => [
                'settings' => [
                    'number_of_shards' => 1,
                    'number_of_replicas' => 1
                ],
                'mappings' => [
                    'properties' => [
                        '@timestamps' => [
                            'type' => 'date',
                            'format' => 'yyyy-MM-dd HH:mm:ss'
                        ],
                        'ip' => ['type' => 'ip'],
                        'url' => ['type' => 'keyword'],
                        'client' => ['type' => 'text']
                    ]
                ]
            ]
        ];
        $this->elasticsearchService->createIndexWithMapping($params);
        $document = [
            '@timestamps' => date('Y-m-d H:i:s'),
            'ip' => $ip,
            'url' => $url,
            'client' => $client
        ];

        return $this->elasticsearchService->index($params['index'], $document);
    }
}
