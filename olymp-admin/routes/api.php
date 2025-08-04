<?php

use App\Http\Controllers\api\ParticipantApiController;
use App\Http\Controllers\api\SchoolApiController;
use App\Http\Controllers\api\UserApiController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::middleware('auth:sanctum')->get('/user', function (Request $request) {
    return $request->user();
});
Route::get('/user/index/{page?}', [UserApiController::class, 'index'])->name('user-api.index');
Route::get('/user/create', [UserApiController::class, 'create'])->name('user-api.create');
Route::post('/user/store', [UserApiController::class, 'store'])->name('user-api.store');
Route::get('/user/edit/{id}', [UserApiController::class, 'edit'])->name('user-api.edit');
Route::put('/user/update/{id}', [UserApiController::class, 'update'])->name('user-api.update');
Route::delete('/user/delete/{id}', [UserApiController::class, 'delete'])->name('user-api.delete');
Route::get('/user/show/{id}', [UserApiController::class, 'show'])->name('user-api.show');

Route::get('/participant/index/{page?}', [ParticipantApiController::class, 'index'])->name('participant-api.index');
Route::get('/participant/create', [ParticipantApiController::class, 'create'])->name('participant-api.create');
Route::post('/participant/store', [ParticipantApiController::class, 'store'])->name('participant-api.store');
Route::get('/participant/edit/{id}', [ParticipantApiController::class, 'edit'])->name('participant-api.edit');
Route::post('/participant/update/{id}', [ParticipantApiController::class, 'update'])->name('participant-api.update');
Route::delete('/participant/delete/{id}', [ParticipantApiController::class, 'delete'])->name('participant-api.delete');
Route::get('/participant/show/{id}', [ParticipantApiController::class, 'show'])->name('participant-api.show');

Route::get('/school/index/{page?}', [SchoolApiController::class, 'index'])->name('school-api.index');
Route::get('/school/create', [SchoolApiController::class, 'create'])->name('school-api.create');
Route::post('/school/store', [SchoolApiController::class, 'store'])->name('school-api.store');
Route::get('/school/edit/{id}', [SchoolApiController::class, 'edit'])->name('school-api.edit');
Route::post('/school/update/{id}', [SchoolApiController::class, 'update'])->name('school-api.update');
Route::delete('/school/delete/{id}', [SchoolApiController::class, 'delete'])->name('school-api.delete');
Route::get('/school/show/{id}', [SchoolApiController::class, 'show'])->name('school-api.show');
