<?php

namespace App\Http\Controllers\api;

use App\Http\Controllers\Controller;
use App\Repositories\UserRepository;
use App\Services\UserService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class UserApiController extends Controller
{
    private UserRepository $userRepository;
    private UserService $userService;
    public function __construct(
        UserRepository $userRepository,
        UserService $userService
    )
    {
        $this->userRepository = $userRepository;
        $this->userService = $userService;
    }

    public function index(Request $request, $page = 1){
        $token =  $request->header('Authorization');
        $usersAmount = $this->userRepository->getCount($token);
        $users = $this->userService->findAll($page, $token);
        return response()->json([
            'page' => $page,
            'users' => array_map(fn($user) => (array)$user, $users),
            'usersAmount' => $usersAmount,
            'currentPage' => $page,
            'perPage' => 10
        ]);
    }
}
