<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use App\Services\TokenService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class ParticipantRepository
{
    private ApiService $apiService;
    private TokenService $tokenService;
    public function __construct(
        ApiService $apiService,
        TokenService $tokenService
    ){
        $this->apiService = $apiService;
        $this->tokenService = $tokenService;
    }
    public function getByApiAll($page = 1, $limit = 10)
    {
        $token = $this->tokenService->getToken();
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_URL_API,
            [
                'page' => $page,
                'limit' => $limit
            ],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return isset($response['data']['data']) ? $response['data']['data'] : [];
    }
    public function getByApiId($id)
    {
        $token = $this->tokenService->getToken();
        $response =  $this->apiService->get(
            ApiHelper::PARTICIPANT_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return isset($response['data']['data']) ? $response['data']['data'] : [];
    }
    public function getByApiUserId($id)
    {
        $token = $this->tokenService->getToken();
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_BY_USER_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return isset($response['data']['data']) ? $response['data']['data'] : [];
    }
    public function getCount()
    {
        $token = $this->tokenService->getToken();
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_COUNT_URL_API,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return isset($response['data']['data']) ? $response['data']['data'] : 0;
    }
}
