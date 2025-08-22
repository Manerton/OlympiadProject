<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Support\Facades\Cookie;

class TokenRepository
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }
    public function revoke($userId){
        $token = request()->header('Authorization');
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $this->apiService->post(ApiHelper::TOKEN_REVOKE_URL_API,  [
            'id' => $userId
        ],
        [
            'Authorization' => "Bearer ". $token
        ]);
    }
    public function refresh()
    {
        return $this->apiService->post(ApiHelper::TOKEN_REFRESH_URL_API, [], [], true);
    }
}
