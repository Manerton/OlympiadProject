<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Support\Facades\Cookie;

class EventJuryRepository
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }

    public function getByEventId($eventId, $token = null){
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::EVENT_JURY_URL_API . '/' . $eventId,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
}
