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
        $token = Request::header('Authorization');
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
