<?php

namespace App\Http\Controllers\api;

use App\Components\Dictionaries\GenderDictionary;
use App\Components\Dictionaries\RoleDictionary;
use App\Components\RabbitMQHelper;
use App\Http\Controllers\Controller;
use App\Http\Requests\UserRequest;
use App\Repositories\TokenRepository;
use App\Repositories\UserRepository;
use App\Services\RabbitMQService;
use App\Services\UserService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class UserApiController extends Controller
{
    private UserRepository $userRepository;
    private RabbitMQService $rabbitMQService;
    private UserService $userService;
    private TokenRepository $tokenRepository;
    public function __construct(
        UserRepository $userRepository,
        RabbitMQService $rabbitMQService,
        UserService $userService,
        TokenRepository $tokenRepository
    )
    {
        $this->userRepository = $userRepository;
        $this->rabbitMQService = $rabbitMQService;
        $this->userService = $userService;
        $this->tokenRepository = $tokenRepository;
    }

    public function index($page = 1){
        $usersAmount = $this->userRepository->getCount();
        $users = $this->userService->findAll($page);
        return response()->json([
            'page' => $page,
            'users' => array_map(fn($user) => (array)$user, $users),
            'usersAmount' => $usersAmount,
            'currentPage' => $page,
            'perPage' => 10
        ]);
    }
    public function create(){
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return response()->json([
           'roles' => $roles,
           'genders' => $genders
        ]);
    }
    public function store(UserRequest $request){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::CREATE,
            RabbitMQHelper::USER_TABLE,
            array_diff_key($data, ['id' => null]),
        );
        return response()->json([
            'status' => 'OK!'
        ]);
    }
    public function show($id){
        $model = $this->userService->find($id);
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return response()->json([
            'model' => (array)$model,
            'roles' => $roles,
            'genders' => $genders
        ]);
    }
    public function edit($id){
        $user = $this->userService->find($id);
        $roles = RoleDictionary::getList();
        $genders = GenderDictionary::getList();
        return response()->json([
            'user' => (array)$user,
            'roles' => $roles,
            'genders' => $genders
        ]);
    }
    public function update(UserRequest $request, $id){
        $data = $request->validated();
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::UPDATE,
            RabbitMQHelper::USER_TABLE,
            array_diff_key($data, ['id' => null]),
            ['id' => $id]
        );
        return response()->json([]);
    }
    public function delete($id){
        $this->rabbitMQService->publish(
            [RabbitMQHelper::AUTH_QUEUE_NAME],
            RabbitMQHelper::QUEUE_NAME,
            RabbitMQHelper::DELETE,
            RabbitMQHelper::USER_TABLE,
            [],
            ['id' => $id]
        );
        return response()->json([]);
    }
    public function revoke($id)
    {
        $this->tokenRepository->revoke($id);
        return response()->json([]);
    }
}
