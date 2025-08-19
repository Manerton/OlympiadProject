<?php

namespace App\Services;

use Elastic\Elasticsearch\Client;
use Elastic\Elasticsearch\ClientBuilder;
use Elastic\Elasticsearch\Exception\AuthenticationException;

class ElasticsearchService
{
    protected Client $client;

    public function __construct()
    {
        try {
            $this->client = ClientBuilder::create()
                ->setHosts([
                    sprintf(
                        '%s://%s:%s',
                        env('ELASTICSEARCH_SCHEME', 'http'),
                        env('ELASTICSEARCH_HOST', 'elasticsearch'),
                        env('ELASTICSEARCH_PORT', '9200')
                    )
                ])
                ->build();
        } catch (\Exception $e) {
            throw new \RuntimeException('Elasticsearch connection failed: ' . $e->getMessage());
        }
    }
    public function createIndexWithMapping(array $params)
    {
        $index = $params['index'];
        $exists = $this->client->indices()->exists(['index' => $index]);
        if ($exists->getStatusCode() != 200) {
            return $this->client->indices()->create($params);
        }
        return true;
    }
    public function getClient(): Client
    {
        return $this->client;
    }

    public function index(string $index, array $document, string $id = null)
    {
        $params = [
            'index' => $index,
            'body'  => $document,
        ];

        if ($id) {
            $params['id'] = $id;
        }

        return $this->client->index($params);
    }

    public function search(string $index, array $query)
    {
        return $this->client->search([
            'index' => $index,
            'body'  => $query,
        ]);
    }
    public function get(string $index, string $id, array $sourceFields = [])
    {
        $params = [
            'index' => $index,
            'id'    => $id,
        ];

        if (!empty($sourceFields)) {
            $params['_source'] = $sourceFields;
        }

        return $this->client->get($params);
    }
    public function deleteIndex(string $index): bool
    {
        try {
            $this->client->indices()->delete([
                'index' => $index
            ]);

            return true; // Успех
        } catch (\Elastic\Elasticsearch\Exception\ClientResponseException $e) {
            if ($e->getCode() === 404) {
                // Индекс не найден
                return false;
            }
            throw $e; // пробрасываем остальные ошибки
        }
    }
    public function decode($response){
        $result = $response->asArray();
        return $result['hits']['hits'];
    }
}
