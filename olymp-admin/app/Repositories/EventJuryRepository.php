<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use App\Services\TokenService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class EventJuryRepository
{
    private ApiService $apiService;
    private TokenService $tokenService;
    public function __construct(
        ApiService $apiService,
        TokenService $tokenService
    )
    {
        $this->apiService = $apiService;
        $this->tokenService = $tokenService;
    }

    public function getByEventId($eventId){
        $token = $this->tokenService->getToken();
        $response = $this->apiService->get(
            ApiHelper::EVENT_JURY_URL_API . '/' . $eventId,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return isset($response['data']['data']) ? $response['data']['data'] : [];
    }
}
