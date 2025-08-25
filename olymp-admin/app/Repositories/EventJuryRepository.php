<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Http\Request;
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

    public function getByEventId($eventId){
        $token = null;
        if (isset(Cookie::getQueuedCookies()[0])) {
            $queuedCookie = Cookie::getQueuedCookies()[0];
            $token = json_decode($queuedCookie->getValue(), true)['token'] ?? null;
        }
        if (is_null($token) && Cookie::get('username')) {
            $token = json_decode(Cookie::get('username'), true)['token'] ?? null;
        }
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
