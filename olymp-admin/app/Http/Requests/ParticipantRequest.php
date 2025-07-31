<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

class ParticipantRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize()
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array
     */
    public function rules()
    {
        return [
            'email' => 'required|email|max:255',
            'firstname' => 'required|string|max:100',
            'password' => 'nullable|string',
            'surname' => 'required|string|max:100',
            'patronymic' => 'nullable|string|max:100',
            'phone_number' => 'required|string',
            'gender' => 'required',
            'role' => 'required',
            'birthdate' => 'required|date|before_or_equal:today',

            'disability' => 'required',
            'citizenship' => 'required',
            'school_id' => 'required',
            'class_number' => 'required'
        ];
    }
}
